from models import db

class UserProfile(db.Model):
    __tablename__ = 'user_profiles'
    
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    first_name = db.Column(db.String(50), nullable=False, default="")
    last_name = db.Column(db.String(50), nullable=False, default="")
    bio = db.Column(db.Text, default="")
    profile_picture = db.Column(db.String(255), default="")
    
    user = db.relationship('User', back_populates='profile')