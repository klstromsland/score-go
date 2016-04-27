class CreateTallies < ActiveRecord::Migration
  def change
    create_table :tallies do |t|
			t.belongs_to :event, index: true
      t.string  :total_time_m, default: "00"
      t.string  :total_time_s, default: "00"
      t.string  :total_time_ms, default: "00"
      t.integer :entrant_tly_id
			t.integer :total_faults, default: 0
      t.decimal :total_points, precision: 5, scale: 2
      t.string  :title
			t.integer	:qualifying_score, default: 0
			t.integer :qualifying_scores, default: 0

      t.timestamps null: false
    end
  end
end