class User < ActiveRecord::Base
  has_and_belongs_to_many :events, inverse_of: :user
  has_secure_password
  def editor?
    self.role == 'editor'
  end
  def admin?
    self.role == 'admin'
  end
  
  def name_role
    "#{last_name}, #{first_name}: #{role}"
  end
	validates :first_name, presence: true
	validates :last_name, presence: true
	validates :role, presence: true
	validates :status, presence: true
	validates :email, presence: true
end