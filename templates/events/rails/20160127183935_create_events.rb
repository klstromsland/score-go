class CreateEvents < ActiveRecord::Migration
  def change
    create_table :events do |t|
      t.string  :title
			t.string  :division
			t.integer :int_search_areas, default: 0
			t.integer :veh_search_areas, default: 0
			t.integer :ext_search_areas, default: 0
			t.integer :cont_search_areas, default: 0
			t.integer :elite_search_areas, default: 0
			t.integer :int_hides, default: 0
			t.integer :veh_hides, default: 0
			t.integer :ext_hides, default: 0
			t.integer :cont_hides, default: 0
			t.integer :elite_hides , default: 0
			t.string  :place
			t.date    :date
			t.string  :host
			t.string  :status

      t.timestamps null: false
    end
  end
end
