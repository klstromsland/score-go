class EventsController < ApplicationController
  before_action :require_user, only: [:index, :edit, :update, :destroy, :create]

	# GET /events
  # GET /events.json
  def index
    @events = Event.all
  end

  def info
  end
  
  # GET /events/1
  # GET /events/1.json
  def show
		@event = Event.find(params[:id])
  end

  # GET /events/new
  def new
    @event = Event.new
  end

  # GET /events/1/edit
  def edit
		@event = Event.find(params[:id])
  end

  # POST /events
  # POST /events.json
  def create
		@event = Event.new(event_params)
		@entrants = Entrant.where(:id => params[:team]) 
		@event.entrants << @entrants
		num_sc_cards = 0
		# For each entrant selected through "team", create a scorecard for each search area
		@event.entrants.each do |entrant|
			Event::ELEMENT_TYPES.each do |elmt|
				case elmt
					when "Container"
						num_sc_cards = @event.cont_search_areas
					when "Exterior"
						num_sc_cards = @event.ext_search_areas
					when "Interior"
						num_sc_cards = @event.int_search_areas
					when "Vehicle"
						num_sc_cards = @event.veh_search_areas
					when "Elite"
            num_sc_cards = @event.elite_search_areas
				end
				if (num_sc_cards != nil)
					for nums in 1..num_sc_cards
						str_tmp = entrant.first_name + " " + entrant.last_name
						@scorecard = Scorecard.new(name: str_tmp, element: elmt, search_area: nums, entrant_scd_id: entrant.id)
						@event.scorecards << @scorecard
					end
				end
			end
			# For each entrant selected through "team", create a tally
			@tally = Tally.new(entrant_tly_id: entrant.id)
			@event.tallies << @tally
		end
		# For each user selected through "user_team", create a user
		@users = User.where(:id => params[:user_team]) 
		@event.users << @users
		respond_to do |format|
			if @event.save
				format.html { redirect_to event_path(@event), notice: 'Event was successfully created.' }
				format.json { render :show, status: :created, location: @event }
			else
				format.html { render :new }
				format.json { render json: @eventinclud.errors, status: :unprocessable_entity }
			end
		end
  end

  # PATCH/PUT /events/1
  # PATCH/PUT /events/1.json
  def update
    @event = Event.find(params[:id])
    @entrants = Entrant.where(:id => params[:team])
    @users = User.where(:id => params[:user_team])
    @event.update_event(@event.id, @entrants, @users)
    respond_to do |format|
      if @event.update(event_params)
        format.html { redirect_to edit_event_path(@event), notice: 'Event was successfully updated.' }
        format.json { render :show, status: :ok, location: @event }
      else
        format.html { render :edit }
        format.json { render json: @event.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /events/1
  # DELETE /events/1.json
  def destroy
		@event = Event.find(params[:id])
		@event.scorecards.destroy_all
		@event.tallies.destroy_all	
		@event.entrants.destroy_all
    @event.destroy
    respond_to do |format|
      format.html { redirect_to events_url, notice: 'Event was successfully destroyed.' }
      format.json { head :no_content }
    end
  end
  private
    # Use callbacks to share common setup or constraints between actions.
#    def set_event
#      @event = Event.find(params[:id])
#    end

    # Never trust parameters from the scary internet, only allow the white list through.

    def event_params
      params.require(:event).permit( :title, :place, :division, :date, :host, :int_search_areas, :ext_search_areas, :cont_search_areas, :veh_search_areas, :elite_search_areas, :int_hides, :ext_hides, :cont_hides, :veh_hides, :elite_hides, :status, :entrant_ids => [], :scorecard_ids => [], :tally_ids => [], :user_ids => [] )
    end
end