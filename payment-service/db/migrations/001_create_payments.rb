Sequel.migration do
  change do
    create_table(:payments) do
      column :order_id, "uuid", primary_key: true
      column :customer_id, "uuid", null: false
      Numeric :amount, size: [12, 2], null: false
      String :status, null: false
      DateTime :updated_at, null: false
    end
  end
end