module PaymentService
  module Domain
    class PaymentReservation
      STATUSES = %w[RESERVED FAILED REFUNDED].freeze

      attr_reader :order_id, :customer_id, :amount, :status

      def initialize(order_id:, customer_id:, amount:, status:)
        @order_id = order_id
        @customer_id = customer_id
        @amount = amount
        @status = status
      end

      def self.reserve(order_id:, customer_id:, amount:)
        new(order_id: order_id, customer_id: customer_id, amount: amount, status: "RESERVED")
      end

      def self.failed(order_id:, customer_id:, amount:)
        new(order_id: order_id, customer_id: customer_id, amount: amount, status: "FAILED")
      end

      def reserved?
        status == "RESERVED"
      end
    end
  end
end