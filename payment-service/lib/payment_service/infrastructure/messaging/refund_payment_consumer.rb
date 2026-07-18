require "racecar"
require "json"
require_relative "../persistence/database"
require_relative "../persistence/sequel_payment_repository"
require_relative "kafka_producer"
require_relative "../../application/use_case/refund_payment"

module PaymentService
  module Infrastructure
    module Messaging
      class RefundPaymentConsumer < Racecar::Consumer
        subscribes_to "payment.refund.command"

        def initialize
          db = PaymentService::Infrastructure::Persistence.connect
          repository = PaymentService::Infrastructure::Persistence::SequelPaymentRepository.new(db)
          @use_case = PaymentService::Application::UseCase::RefundPayment.new(repository: repository)
          @producer = KafkaProducer.new(brokers: ENV.fetch("KAFKA_BROKERS", "localhost:9092"))
        end

        def process(message)
          command = JSON.parse(message.value)
          payload = command["payload"]

          success = @use_case.execute(order_id: payload["orderId"])

          reply = {
            replyType: "payment.refund.reply",
            sagaId: command["sagaId"],
            success: success,
            occurredAt: Time.now.utc.iso8601
          }

          @producer.publish(topic: "payment.refund.reply", key: command["sagaId"], payload: reply.to_json)
          puts "[payment.refund.command] order #{payload['orderId']} refunded=#{success}"
        rescue => e
          puts "[payment.refund.command] failed: #{e.message}"
          raise
        end
      end
    end
  end
end