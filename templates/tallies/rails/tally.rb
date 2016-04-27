class Tally < ActiveRecord::Base
  belongs_to :event

  def get_tally(tally_id)
    current_tally = Tally.find_by(id: tally_id)
    specialty = ""
    q_scores = 0
    point_tally = 0.0
    if current_tally.event.division == "Element Specialty"
      if current_tally.event.cont_search_areas > 0
        specialty = "Container"
      elsif current_tally.event.ext_search_areas > 0
        specialty = "Exterior"
      elsif current_tally.event.int_search_areas > 0
        specialty = "Interior"
      else
        specialty = "Vehicle"
      end
      current_tally.event.entrants.each do |entrant|
        if current_tally.entrant_tly_id == entrant.id
          entrant.events.each do |event|
            if event.id != current_tally.event.id
              if event.division == "Element Specialty"
                if specialty == "Container" && event.cont_search_areas > 0 || specialty == "Exterior" && event.ext_search_areas > 0 || specialty == "Interior" && event.int_search_areas > 0 || specialty == "Vehicle" && event.veh_search_areas > 0
                  event.tallies.each do |tally|
                    if tally.entrant_tly_id == entrant.id  
                      if tally.qualifying_score == 1
                        q_scores += 1
                      end                      
                    end
                  end
                end
              end
            end
          end  
        end
      end        
    end
    if q_scores >= 2
      q_scores = 0
    end
    if current_tally.event.division == "Elite"
      current_tally.event.entrants.each do |entrant|
        if current_tally.entrant_tly_id == entrant.id
          entrant.events.each do |event|
            if event.id != current_tally.event.id
              if event.division == "Elite"
                event.tallies.each do |tally|
                  if tally.entrant_tly_id == entrant.id                
                    if tally.total_points != nil
                      point_tally += tally.total_points
                    end
                  end
                end
              end
            end
          end
        end
      end
    end    
    current_tally.event.entrants.each do |entrant|
      if current_tally.entrant_tly_id == entrant.id
#        point_tally = 0.0
        time_tally_m = 0
        time_tally_mstr = ""
        time_tally_s = 0
        time_tally_sstr = ""
        time_tally_ms = 0
        time_tally_msstr = ""
        fault_tally = 0
        q_score = 0
#        q_scores = 0
        titled = " "
        current_tally.event.scorecards.each do |scorecard|
          if scorecard.entrant_scd_id == entrant.id
            if scorecard.total_points != nil
              point_tally += scorecard.total_points
            end
            if scorecard.time_elapsed_m != nil
                time_tally_m += scorecard.time_elapsed_m.to_i
            end
            if scorecard.time_elapsed_s != nil	
                time_tally_s += scorecard.time_elapsed_s.to_i 
            end
            if scorecard.time_elapsed_ms != nil
                time_tally_ms += scorecard.time_elapsed_ms.to_i 
            end
            if scorecard.total_faults != nil
                fault_tally += scorecard.total_faults
            end            
            if point_tally.round == 100 && fault_tally <= 3 && current_tally.event.division != "Element Specialty" && current_tally.event.division != "Elite"
              titled = "Titled"	
            elsif current_tally.event.division == "Element Specialty"
              if point_tally.round >= 75 && fault_tally <= 3
                q_score = 1
                q_scores += 1
                if q_scores == 2
                  titled = "Titled"
                end
              end
              if point_tally.round == 100 && fault_tally <= 3
                titled = "Titled"
              end
            elsif current_tally.event.division == "Elite" && point_tally.round >= 150 && fault_tally <= 3
              titled = "Titled"
            else
              titled = "Not this time"
            end
            time_tally_mstr = time_tally_m.to_s
            time_tally_sstr = time_tally_s.to_s
            time_tally_msstr = time_tally_ms.to_s
            if time_tally_m < 10	    
              time_tally_mstr = "0" + time_tally_mstr
            end
            if time_tally_s < 10
              time_tally_sstr = "0" + time_tally_sstr
            end
            if time_tally_ms < 10
              time_tally_msstr = "0" + time_tally_msstr
            end
            update_attribute(:total_time_m, time_tally_mstr)
            update_attribute(:total_time_s, time_tally_sstr)
            update_attribute(:total_time_ms, time_tally_msstr)
            update_attribute(:total_points, point_tally)
            update_attribute(:total_faults, fault_tally)
            update_attribute(:title, titled)         
            update_attribute(:qualifying_score, q_score)
            update_attribute(:qualifying_scores, q_scores)            
          end
        end	
      end
    end
	end
 
  def time_to_ms(tally_id)
    current_tally = Tally.find_by(id: tally_id)
    return current_tally.total_time_m.to_i * 60000 + current_tally.total_time_s.to_i * 1000 + current_tally.total_time_ms.to_i * 10
  end 
end