from models import db

class UserProfile(db.Model):
    __tablename__ = 'user_profile'
    
    id = db.Column(db.Integer, primary_key = True)
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable = False)
    first_name = db.Column(db.String(80), nullable = False)
    last_name = db.Column(db.String(80), nullable = False)
    bio = db.Column(db.Text, nullable = True)
    profile_picture = db.Column(db.String(255), nullable = True)
    
    user = db.relationship('User', back_populates='profile')
    