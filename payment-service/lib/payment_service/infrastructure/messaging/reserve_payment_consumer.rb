require "racecar"
require "json"
require_relative "../persistence/database"
require_relative "../persistence/sequel_payment_repository"
require_relative "kafka_producer"
require_relative "../../application/use_case/reserve_payment"

module PaymentService
  module Infrastructure
    module Messaging
      class ReservePaymentConsumer < Racecar::Consumer
        subscribes_to "payment.reserve.command"

        def initialize
          db = PaymentService::Infrastructure::Persistence.connect
          repository = PaymentService::Infrastructure::Persistence::SequelPaymentRepository.new(db)
          @use_case = PaymentService::Application::UseCase::ReservePayment.new(repository: repository)
          @producer = KafkaProducer.new(brokers: ENV.fetch("KAFKA_BROKERS", "localhost:9092"))
        end

        def process(message)
          command = JSON.parse(message.value)
          payload = command["payload"]

          reservation = @use_case.execute(
            order_id: payload["orderId"],
            customer_id: payload["customerId"],
            total_amount: payload["totalAmount"].to_f
          )

          reply = {
            replyType: "payment.reserve.reply",
            sagaId: command["sagaId"],
            success: reservation.reserved?,
            reason: reservation.reserved? ? nil : "invalid amount",
            occurredAt: Time.now.utc.iso8601
          }

          @producer.publish(topic: "payment.reserve.reply", key: command["sagaId"], payload: reply.to_json)
          puts "[payment.reserve.command] order #{payload['orderId']} -> #{reservation.status}"
        rescue => e
          puts "[payment.reserve.command] failed: #{e.message}"
          raise
        end
      end
    end
  end
end