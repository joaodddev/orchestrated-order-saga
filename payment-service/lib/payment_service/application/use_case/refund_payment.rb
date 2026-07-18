require_relative "../../domain/payment_reservation"

module PaymentService
  module Application
    module UseCase
      class RefundPayment
        def initialize(repository:)
          @repository = repository
        end

        def execute(order_id:)
          reservation = @repository.find_by_order_id(order_id)
          return false unless reservation

          refunded = PaymentService::Domain::PaymentReservation.new(
            order_id: reservation.order_id, customer_id: reservation.customer_id,
            amount: reservation.amount, status: "REFUNDED"
          )
          @repository.upsert(refunded)
          true
        end
      end
    end
  end
end