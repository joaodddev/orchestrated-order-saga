require "sequel"
require "dotenv/load"

module PaymentService
  module Infrastructure
    module Persistence
      def self.connect
        Sequel.connect(
          adapter: "postgres",
          host: ENV.fetch("DB_HOST", "localhost"),
          port: ENV.fetch("DB_PORT", "5433"),
          user: ENV.fetch("DB_USER", "saga"),
          password: ENV.fetch("DB_PASSWORD", "saga"),
          database: ENV.fetch("DB_NAME", "payment_db")
        )
      end
    end
  end
end