class CreateScorecards < ActiveRecord::Migration
  def change
    create_table :scorecards do |t|
			t.belongs_to :event, index: true
			t.integer :entrant_scd_id
			t.integer :hides_found, default: 0
			t.decimal :hides_missed, default: 0
			t.float   :maxpoint, precision: 5, scale: 2
			t.integer :hides_max, default: 0
			t.integer :other_faults_count, default: 0
			t.integer :search_area
			t.integer :total_faults
			t.decimal :total_points, precision: 5, scale: 2
			t.string  :absent, default: "no"
			t.string  :element
			t.string  :dismissed, default: "no"
			t.string  :eliminated_during_search, default: "no"
			t.string  :excused, default: "no"
			t.integer :false_alert_fringe, default: "no"
			t.string  :finish_call, default: "yes"
			t.string  :judge_signature, default: "no"
			t.string  :maxtime_m, default: "00"
			t.string  :maxtime_s, default: "00"
			t.string  :maxtime_ms, default: "00"
      t.string  :name
			t.string  :pronounced, default: "no"
			t.string  :timed_out
			t.string  :time_elapsed_m, default: "00"
			t.string  :time_elapsed_s, default: "00"  
			t.string  :time_elapsed_ms, default: "00"
			t.text    :comments
			t.text    :other_faults_descr

      t.timestamps null: false
    end
  end
end
