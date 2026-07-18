require_relative "../lib/payment_service/infrastructure/messaging/reserve_payment_consumer"
require_relative "../lib/payment_service/infrastructure/messaging/refund_payment_consumer"

Racecar.configure do |config|
  config.brokers = [ENV.fetch("KAFKA_BROKERS", "localhost:9092")]
  config.client_id = "payment-service"
  config.group_id = "payment-service-group"
end