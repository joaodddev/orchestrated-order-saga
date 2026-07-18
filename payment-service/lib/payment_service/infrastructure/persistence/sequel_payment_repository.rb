require_relative "../../application/port/payment_repository"
require_relative "../../domain/payment_reservation"

module PaymentService
  module Infrastructure
    module Persistence
      class SequelPaymentRepository < PaymentService::Application::Port::PaymentRepository
        def initialize(db)
          @db = db
        end

        def upsert(reservation)
          @db[:payments].insert_conflict(target: :order_id, update: { status: reservation.status, updated_at: Time.now.utc }).insert(
            order_id: reservation.order_id,
            customer_id: reservation.customer_id,
            amount: reservation.amount,
            status: reservation.status,
            updated_at: Time.now.utc
          )
        end

        def find_by_order_id(order_id)
          row = @db[:payments].where(order_id: order_id).first
          return nil unless row

          PaymentService::Domain::PaymentReservation.new(
            order_id: row[:order_id], customer_id: row[:customer_id],
            amount: row[:amount].to_f, status: row[:status]
          )
        end
      end
    end
  end
end