class CreateEntrants < ActiveRecord::Migration
  def change
    create_table :entrants do |t|
      t.string :first_name
      t.string :last_name
      t.string :idnumber
      t.string :dog_name
      t.string :dogidnumber
      t.string :breed

      t.timestamps null: false
    end
  end
end
