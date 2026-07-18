require_relative "../../domain/payment_reservation"

module PaymentService
  module Application
    module UseCase
      class ReservePayment
        def initialize(repository:)
          @repository = repository
        end

        def execute(order_id:, customer_id:, total_amount:)
          reservation = if total_amount.positive?
            PaymentService::Domain::PaymentReservation.reserve(order_id: order_id, customer_id: customer_id, amount: total_amount)
          else
            PaymentService::Domain::PaymentReservation.failed(order_id: order_id, customer_id: customer_id, amount: total_amount)
          end

          @repository.upsert(reservation)
          reservation
        end
      end
    end
  end
end