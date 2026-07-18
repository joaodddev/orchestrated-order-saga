module PaymentService
  module Application
    module Port
      class PaymentRepository
        def upsert(reservation)
          raise NotImplementedError
        end

        def find_by_order_id(order_id)
          raise NotImplementedError
        end
      end
    end
  end
end