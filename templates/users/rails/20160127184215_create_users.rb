class CreateUsers < ActiveRecord::Migration
  def change
    create_table :users do |t|
      t.string :first_name
      t.string :last_name
      t.string :role
      t.string :approved
      t.string :status
			t.string :email
			t.string :password_digest

      t.timestamps null: false
    end
  end
end
