class Event < ActiveRecord::Base
  ELEMENT_TYPES = ["Container", "Exterior", "Interior", "Vehicle", "Elite"]
  has_and_belongs_to_many :entrants, inverse_of: :event
  has_many :scorecards, dependent: :destroy
  has_many :tallies, dependent: :destroy
  has_and_belongs_to_many :users, inverse_of: :event
	validates :title, presence: true
	validates :division, presence: true
	validates :host, presence: true
	validates :int_search_areas, numericality:
	{greater_than_or_equal_to: 0}
	validates :ext_search_areas, numericality:
	{greater_than_or_equal_to: 0}
	validates :cont_search_areas, numericality:
	{greater_than_or_equal_to: 0}
	validates :veh_search_areas, numericality:
	{greater_than_or_equal_to: 0}
	validates :elite_search_areas, numericality:
	{greater_than_or_equal_to: 0}
  validates :int_hides, numericality:
	{greater_than_or_equal_to: 0}
	validates :ext_hides, numericality:
	{greater_than_or_equal_to: 0}
	validates :cont_hides, numericality:
	{greater_than_or_equal_to: 0}
	validates :veh_hides, numericality:
	{greater_than_or_equal_to: 0}
	validates :elite_hides, numericality:
	{greater_than_or_equal_to: 0}

  def scorecard_completion(event_id)
    current_event = Event.find_by(id: event_id)
    completed_entrant_scorecards = Array.new(current_event.entrants.length, 0)
    count = 0
    entrant_count = 0
    completed_entrant_scorecards.each do |completed|
      if entrant_count >= current_event.entrants.length
        break
      end
      current_event.entrants.each do |entrant|
        sc_count = 0
        Event::ELEMENT_TYPES.each do |elmt|
          case elmt
            when "Container"
              sa_count = current_event.cont_search_areas 
            when "Exterior"
              sa_count = current_event.ext_search_areas
            when "Interior"
              sa_count = current_event.int_search_areas
            when "Vehicle"
              sa_count = current_event.veh_search_areas
            when "Elite"
              sa_count = current_event.elite_search_areas			
          end
          current_event.scorecards.each do |scorecard|
            if (scorecard.entrant_scd_id == entrant.id) && (scorecard.element == elmt)
              if scorecard.judge_signature == "yes"
                sc_count += 1
              end
            end
          end
        end
        if sc_count == current_event.cont_search_areas + current_event.ext_search_areas + current_event.int_search_areas + current_event.veh_search_areas + current_event.elite_search_areas
          completed_entrant_scorecards[count] = entrant.id
          count += 1
        end
      end
      entrant_count += 1
    end
    return completed_entrant_scorecards
  end

  def tally_completion(event_id)
    current_event = Event.find_by(id: event_id)
    completed_entrant_tallies = Array.new(current_event.entrants.length, 0)
    count = 0
    entrant_count = 0
    completed_entrant_tallies.each do |completed|
      if entrant_count >= current_event.entrants.length
        break
      end    
      current_event.entrants.each do |entrant|
        current_event.tallies.each do |tally|
          if tally.entrant_tly_id == entrant.id
            if tally.total_points != nil
              completed_entrant_tallies[count] = entrant.id
              count += 1
            end
          end
        end
      end
      entrant_count += 1
    end
    return completed_entrant_tallies
  end
  
  def update_event(event_id, new_entrants, new_users)
    current_event = Event.find_by(id: event_id)
    if current_event != nil    
      sc_count = 0
      entrant_counter = 0
      sa_count = 0
      nums = 0
      current_event.update_attribute(:cont_search_areas, current_event.cont_search_areas)
      current_event.update_attribute(:ext_search_areas, current_event.ext_search_areas)
      current_event.update_attribute(:int_search_areas, current_event.int_search_areas)
      current_event.update_attribute(:veh_search_areas, current_event.veh_search_areas)
      current_event.update_attribute(:elite_search_areas, current_event.elite_search_areas)
      current_event.update_attribute(:cont_hides, current_event.cont_hides)    
      current_event.update_attribute(:ext_hides, current_event.ext_hides)
      current_event.update_attribute(:int_hides, current_event.int_hides)
      current_event.update_attribute(:veh_hides, current_event.veh_hides)
      current_event.update_attribute(:elite_hides, current_event.elite_hides)    
      # For each element, check if changes have been made to search area entries for each entrant and update if needed.
      Event::ELEMENT_TYPES.each do |elmt|
        case elmt
          when "Container"
            sa_count = current_event.cont_search_areas 
          when "Exterior"
            sa_count = current_event.ext_search_areas
          when "Interior"
            sa_count = current_event.int_search_areas
          when "Vehicle"
            sa_count = current_event.veh_search_areas
          when "Elite"
            sa_count = current_event.elite_search_areas			
        end
        # For element search area, check if changes have been made to search area entries for each entrant and update if needed.
        sc_count = 0
        current_event.entrants.each do |entrant|
          sc_count = 0
          # Count number of scorecards for element for entrant.
          current_event.scorecards.each do |scorecard|
            if (scorecard.entrant_scd_id == entrant.id) && (scorecard.element == elmt)
              sc_count += 1
            end
          end
          # For entrant, if scorecard count equals element search area count, no change. Go to next entrant for this search area.
          if sc_count == sa_count
            next
          end
          # For entrant, if scorecard count does not equal element search area count, update.
          # If scorecard count = 0, add first scorecard, increment count, check condition
          if (sc_count == 0) && (sa_count > sc_count)
            str_tmp = entrant.first_name + " " + entrant.last_name
            @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: 1, entrant_scd_id: entrant.id)
            current_event.scorecards << @scorecard
            sc_count += 1
            if sc_count == sa_count
              break
            end
          end
          if (sc_count > 0) && (sa_count > sc_count)
            # If scorecard count is greater than zero, find after which index (scorecard.search_area)	to add scorecard and add it, increment count, check condition			
            current_event.scorecards.each do |scorecard|
              if scorecard.entrant_scd_id == entrant.id
                str_tmp = entrant.first_name + " " + entrant.last_name
                @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: sc_count + 1, entrant_scd_id: entrant.id)
                current_event.scorecards << @scorecard
                sc_count += 1
              end
              if sc_count == sa_count
                break
              end
            end
          end
          # If scorecard count is greater than search area count, find a scorecard to delete with index greater than the number of search areas and delete it, deprecate count, check condition.
          if sa_count < sc_count 
            current_event.scorecards.each do |scorecard|
              if (scorecard.entrant_scd_id == entrant.id) && (scorecard.search_area > sa_count)
                current_event.scorecards.destroy(scorecard)
                sc_count -= 1
              end
              if sc_count == sa_count
                break
              end
            end
          end
        end
      end		
      #Check if changes have been made to entrants selected through "team" and update if needed.		
      #Check if entrant count has changed.
      if (new_entrants.length == 0) && (current_event.entrants.length > 0)
        current_event.entrants.each do |eventrant|
          current_event.scorecards.each do |scorecard|
            if scorecard.entrant_scd_id == eventrant.id
              current_event.scorecards.destroy(scorecard)
            end
          end
          current_event.tallies.each do |tally|
            if tally.entrant_tly_id == eventrant.id
              current_event.tallies.destroy(tally)
            end
          end
          current_event.entrants.destroy(eventrant)
        end
      elsif (current_event.entrants.length == 0) && (new_entrants.length > 0)
        new_entrants.each do |entrant|
          Event::ELEMENT_TYPES.each do |elmt|
            case elmt
              when "Container"
                sa_count = current_event.cont_search_areas
              when "Exterior"
                sa_count = current_event.ext_search_areas
              when "Interior"
                sa_count = current_event.int_search_areas
              when "Vehicle"
                sa_count = current_event.veh_search_areas
              when "Elite"
                sa_count = current_event.elite_search_areas			
            end
            if sa_count >= 0	  
              for nums in 1..sa_count
                str_tmp = entrant.first_name + " " + entrant.last_name
                @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: nums, entrant_scd_id: entrant.id)
                current_event.scorecards << @scorecard
              end
            end
          end
          @tally = Tally.new(entrant_tly_id: entrant.id)
          current_event.tallies << @tally
          current_event.entrants << entrant
        end
      elsif (current_event.entrants.length > new_entrants.length) && (current_event.entrants.length != 0)
        entrant_counter = 0
        current_event_entrants_tmp = current_event.entrants
        current_event_entrants_tmp.each do |eventrant|	
          if entrant_counter < new_entrants.length
            new_entrants.each do |entrant|
              if entrant.id == eventrant.id
                entrant_counter += 1
                break
              else
                current_event.scorecards.each do |scorecard|
                  if scorecard.entrant_scd_id == eventrant.id
                    current_event.scorecards.destroy(scorecard)
                  end
                end
                current_event.tallies.each do |tally|
                  if tally.entrant_tly_id == eventrant.id
                    current_event.tallies.destroy(tally)
                  end
                end
                current_event.entrants.destroy(eventrant)						
              end
            end
            if current_event.entrants.length == new_entrants.length
              break
            end
          else
            current_event.scorecards.each do |scorecard|
              if scorecard.entrant_scd_id == eventrant.id
                current_event.scorecards.destroy(scorecard)
              end
            end
            current_event.tallies.each do |tally|
              if tally.entrant_tly_id == eventrant.id
                current_event.tallies.destroy(tally)
              end
            end
            current_event.entrants.destroy(eventrant)
            if current_event.entrants.length == new_entrants.length
              break
            end					
          end
        end
      elsif (current_event.entrants.length < new_entrants.length) && (new_entrants.length != 0)
        entrant_counter = 0
        new_entrant_counter = 0
        current_event_entrants_tmp = current_event.entrants
        #find existing entrants that match new entrants so that can be ignored?
        new_entrants.each do |entrant|				
          if entrant_counter < current_event_entrants_tmp.length
            current_event_entrants_tmp.each do |eventrant|			
              if entrant.id == eventrant.id
                entrant_counter += 1
                new_entrant_counter = 0
                break
              elsif entrant_counter < current_event_entrants_tmp.length
                new_entrant_counter += 1
                if new_entrant_counter < current_event_entrants_tmp.length
                  next
                elsif new_entrant_counter == current_event_entrants_tmp.length
                  Event::ELEMENT_TYPES.each do |elmt|
                    case elmt
                      when "Container"
                        sa_count = current_event.cont_search_areas
                      when "Exterior"
                        sa_count = current_event.ext_search_areas
                      when "Interior"
                        sa_count = current_event.int_search_areas
                      when "Vehicle"
                        sa_count = current_event.veh_search_areas
                      when "Elite"
                        sa_count = current_event.elite_search_areas			
                    end
                    if sa_count >= 0	  
                      for nums in 1..sa_count
                        str_tmp = entrant.first_name + " " + entrant.last_name
                        @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: nums, entrant_scd_id: entrant.id)
                        current_event.scorecards << @scorecard
                      end
                    end
                  end
                  @tally = Tally.new(entrant_tly_id: entrant.id)
                  current_event.tallies << @tally
                  current_event.entrants << entrant
                  new_entrant_counter = 0
                end
              end
            end
          elsif current_event.entrants.length < new_entrants.length
            Event::ELEMENT_TYPES.each do |elmt|
              case elmt
                when "Container"
                  sa_count = current_event.cont_search_areas
                when "Exterior"
                  sa_count = current_event.ext_search_areas
                when "Interior"
                  sa_count = current_event.int_search_areas
                when "Vehicle"
                  sa_count = current_event.veh_search_areas
                when "Elite"
                  sa_count = current_event.elite_search_areas			
              end
              if sa_count >= 0	  
                for nums in 1..sa_count
                  str_tmp = entrant.first_name + " " + entrant.last_name
                  @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: nums, entrant_scd_id: entrant.id)
                  current_event.scorecards << @scorecard
                end
              end
            end
            @tally = Tally.new(entrant_tly_id: entrant.id)
            current_event.tallies << @tally
            current_event.entrants << entrant
            if current_event.entrants.length == new_entrants.length
              break
            end
          end
        end
      elsif (current_event.entrants.length == new_entrants.length) && (new_entrants.length != 0)
      # check if entrant ids have changed
        entrant_counter = 0
        current_event_entrants_tmp = current_event.entrants
        current_event_entrants_tmp.each do |eventrant|
          if entrant_counter == current_event_entrants_tmp.length
            break
          end
          new_entrants.each do |entrant|			
            if entrant.id == eventrant.id
              entrant_counter += 1
              break
            elsif entrant_counter <= current_event_entrants_tmp.length
              next
            else
              current_event.scorecards.each do |scorecard|
                if scorecard.entrant_scd_id == eventrant.id
                  current_event.scorecards.destroy(scorecard)
                end
              end
              current_event.tallies.each do |tally|
                if tally.entrant_tly_id == eventrant.id
                  current_event.tallies.destroy(tally)
                end
              end
              current_event.entrants.destroy(eventrant)						
              Event::ELEMENT_TYPES.each do |elmt|
                case elmt
                  when "Container"
                    sa_count = current_event.cont_search_areas
                  when "Exterior"
                    sa_count = current_event.ext_search_areas
                  when "Interior"
                    sa_count = current_event.int_search_areas
                  when "Vehicle"
                    sa_count = current_event.veh_search_areas
                  when "Elite"
                    sa_count = current_event.elite_search_areas			
                end
                if sa_count >= 0	  
                  for nums in 1..sa_count
                    str_tmp = entrant.first_name + " " + entrant.last_name
                    @scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: nums, entrant_scd_id: entrant.id)
                    current_event.scorecards << @scorecard
                  end
                end
              end
              @tally = Tally.new(entrant_tly_id: entrant.id)
              current_event.tallies << @tally
              current_event.entrants << entrant
              entrant_counter += 1
              if entrant_counter == current_event_entrants_tmp.length
                break
              end
            end
            if entrant_counter == current_event_entrants_tmp.length
              break
            end
          end
        end
      end
      current_event.users.destroy_all
      current_event.users << new_users
      current_event.save()
    end
  end
  
  
  def place_order(event_id)
    place_points = 0.0
    place_faults = 0
    place_time = 0
    count = 0
    current_event = Event.find_by(id: event_id)
    placing = Array.new(current_event.tallies.length, 0)   
    placing.each do |place|
      place_points = 0.0
      place_faults = 0
      place_time = 0
      current_event.tallies.each do |tally|
        if placing.include?(tally.id)
          next
        else
          if tally.total_points == nil
            tally.total_points = 0
          end
          if place_points < tally.total_points
            place_time = tally.time_to_ms(tally.id)
            place_points = tally.total_points
            place_faults = tally.total_faults
          end          
          if place_points <= tally.total_points && place_time >= tally.time_to_ms(tally.id) && place_faults >= tally.total_faults
            place_time = tally.time_to_ms(tally.id)
            place_points = tally.total_points
            place_faults = tally.total_faults
          end
        end
      end
      current_event.tallies.each do |tally|
        if placing.include?(tally.id)
          next
        elsif (tally.total_points == place_points) && (tally.time_to_ms(tally.id) == place_time) && (tally.total_faults == place_faults)
          placing[count] = tally.id
          count += 1
          break
        end
      end  
    end
    return placing
  end
end