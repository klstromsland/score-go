class Scorecard < ActiveRecord::Base
  belongs_to :event	
  validates :maxtime_m, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :maxtime_s, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :maxtime_ms, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :time_elapsed_m, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :time_elapsed_s, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :time_elapsed_ms, format: { with: %r{[0-9]+[0-9]}, message: 'must be 2 numeric values' }
  validates :hides_found, numericality: { less_than_or_equal_to: :hides_max, message: 'must be less than or equal to Hides/Calls' }
  def get_elmSearchAreas(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    case current_scorecard.element
      when "Container"
        return current_scorecard.event.cont_search_areas
      when "Interior"
        return current_scorecard.event.int_search_areas
      when "Exterior"
        return current_scorecard.event.ext_search_areas
      when "Vehicle"
        return current_scorecard.event.veh_search_areas
      when "Elite"
        return current_scorecard.event.elite_search_areas
    end
  end

  def get_elmHides(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    case current_scorecard.element
      when "Container"
        return current_scorecard.event.cont_hides
      when "Interior"
        return current_scorecard.event.int_hides
      when "Exterior"
        return current_scorecard.event.ext_hides
      when "Vehicle"
        return current_scorecard.event.veh_hides
      when "Elite"
        return current_scorecard.event.elite_hides
    end
  end
  
  def get_check_hide_count(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    hideCountCheck = current_scorecard.get_elmHides(scorecard_id)
    elm_hides = hideCountCheck
    message = ""
    if current_scorecard.hides_max != nil
      current_scorecard.event.entrants.each do |entrant|
        if current_scorecard.entrant_scd_id == entrant.id
          current_scorecard.event.scorecards.each do |scorecard|
            if scorecard.entrant_scd_id == entrant.id && scorecard.element == current_scorecard.element && scorecard.hides_max != nil
              hideCountCheck -= scorecard.hides_max
            end
          end
        end
      end     
      if (hideCountCheck > elm_hides || hideCountCheck < 0 ) && (current_scorecard.event.division != "NW1")
        errors.add(:base, 'Incorrect Hide Count...')
        message =  "Incorrect Hide Count..."
      elsif current_scorecard.event.division == "NW1"
        if current_scorecard.hides_max != 1
          errors.add(:base, 'Incorrect Hide Count...')
          message = "Incorrect Hide Count..." 
        end
      elsif current_scorecard.event.division == "NW2"
        if current_scorecard.hides_max == 0
          errors.add(:base, 'Incorrect Hide Count...')
          message =  "Incorrect Hide Count..." 
        end
      else
        message =  ""
        return message
      end
    end
  end
  
  def get_max_point(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    point = 0.0
    if current_scorecard.maxpoint != nil
      if current_scorecard.event.division != "NW1"
        case current_scorecard.element
          when "Container"
            if current_scorecard.hides_max > 0
              if current_scorecard.event.division != "Element Specialty"
                point = (25.00/current_scorecard.event.cont_hides.to_f) * current_scorecard.hides_max.to_f
              elsif current_scorecard.event.division == "Element Specialty"
                point = (100.00/current_scorecard.event.cont_hides.to_f) * current_scorecard.hides_max.to_f
              else
                point = 0
              end
            end
          when "Interior"
            if current_scorecard.hides_max > 0
              if current_scorecard.event.division != "Element Specialty"            
                point = (25.00/current_scorecard.event.int_hides.to_f) * current_scorecard.hides_max.to_f
              elsif current_scorecard.event.division == "Element Specialty"
                point = (100.00/current_scorecard.event.int_hides.to_f) * current_scorecard.hides_max.to_f
              else
                point = 0
              end
            end
          when "Exterior"
            if current_scorecard.hides_max > 0
              if current_scorecard.event.division != "Element Specialty"            
                point = (25.00/current_scorecard.event.ext_hides.to_f) * current_scorecard.hides_max.to_f
              elsif current_scorecard.event.division == "Element Specialty"
                point = (100.00/current_scorecard.event.ext_hides.to_f) * current_scorecard.hides_max.to_f              
              else
                point = 0
              end
            end
          when "Vehicle"
            if current_scorecard.hides_max > 0
              if current_scorecard.event.division != "Element Specialty"
                point = (25.00/current_scorecard.event.veh_hides.to_f) * current_scorecard.hides_max.to_f
              elsif current_scorecard.event.division == "Element Specialty"
                point = (100.00/current_scorecard.event.veh_hides.to_f) * current_scorecard.hides_max.to_f
              else
                point = 0
              end
            end
          when "Elite"
            if current_scorecard.hides_max > 0
              point = (100.00/current_scorecard.event.elite_hides.to_f) * current_scorecard.hides_max.to_f
            else
              point = 0
            end
        end
      else
        point =  25.0
      end
      update_attribute(:maxpoint, point)
    end
  end

  def get_fault_total(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    totalfaults = 0.0
    totalfaults = current_scorecard.other_faults_count
    if current_scorecard.false_alert_fringe > 0
      if current_scorecard.event.division != "Elite"
        totalfaults += 2
      end
    end
    if current_scorecard.eliminated_during_search == "yes" || current_scorecard.excused == "yes"
      if current_scorecard.event.division != "Elite"
        totalfaults += 3
      else
        totalfaults += 1
      end
    end
    if (current_scorecard.absent == "yes")&&(event.division != "Elite")
      totalfaults += 4
    end
    update_attribute(:total_faults, totalfaults)
  end

  def get_time(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    if (current_scorecard.timed_out == "yes" || current_scorecard.finish_call == "no") && current_scorecard.event.division == "Elite"
        update_attribute(:time_elapsed_m, current_scorecard.maxtime_m)
        update_attribute(:time_elapsed_s, current_scorecard.maxtime_s)
        update_attribute(:time_elapsed_ms, current_scorecard.maxtime_ms)
    elsif current_scorecard.event.division != "Elite"
      if current_scorecard.timed_out == "yes" || current_scorecard.finish_call == "no" || current_scorecard.absent == "yes" || current_scorecard.eliminated_during_search == "yes" || current_scorecard.excused == "yes" || current_scorecard.false_alert_fringe > 0
        update_attribute(:time_elapsed_m, current_scorecard.maxtime_m)
        update_attribute(:time_elapsed_s, current_scorecard.maxtime_s)
        update_attribute(:time_elapsed_ms, current_scorecard.maxtime_ms)
      end
    end
  end

  def get_points(scorecard_id)
    current_scorecard = Scorecard.find_by(id: scorecard_id)
    points = 0.0
    if current_scorecard.total_faults != nil && current_scorecard.maxpoint != nil && current_scorecard.hides_max != nil && current_scorecard.total_points != nil && current_scorecard.maxpoint != 0
      if current_scorecard.total_faults <= 3
        if current_scorecard.hides_found == current_scorecard.hides_max && current_scorecard.event.division == "NW1"
          points = current_scorecard.maxpoint
        elsif (current_scorecard.event.division == "NW2")||(current_scorecard.event.division == "NW3") || (current_scorecard.event.division == "Element Specialty")
          points = current_scorecard.hides_found * current_scorecard.maxpoint/current_scorecard.hides_max
        elsif current_scorecard.event.division == "Elite"
          if current_scorecard.false_alert_fringe == 0
            points = (current_scorecard.hides_found * current_scorecard.maxpoint/current_scorecard.hides_max) - current_scorecard.total_faults
          elsif current_scorecard.false_alert_fringe <= 3 && current_scorecard.false_alert_fringe > 0
            points = (current_scorecard.hides_found * current_scorecard.maxpoint/current_scorecard.hides_max) - current_scorecard.total_faults + (current_scorecard.false_alert_fringe * 100.0/current_scorecard.event.elite_hides/2)
          end
          if current_scorecard.finish_call == "no"
            if current_scorecard.hides_found > 0 || current_scorecard.false_alert_fringe > 0
              points -= 100.0/current_scorecard.event.elite_hides/2
            end
          end
        else
          points = 0.0
        end
      else
        points = 0.0
      end
      if (current_scorecard.absent == "yes")||(current_scorecard.eliminated_during_search == "yes")||(current_scorecard.excused == "yes")
        points = 0.0
      end
    end
    update_attribute(:total_points, points)
  end
end