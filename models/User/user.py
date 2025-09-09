from models.UserProfile.user_profile import UserProfile
from werkzeug.security import generate_password_hash, check_password_hash
from models import db
from sqlalchemy import event

class User(db.Model):
    __tablename__ = 'users'
    
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    email = db.Column(db.String(120), unique=True, nullable=False)  # Add email field
    _password = db.Column("password", db.String(255), nullable=False)    
    
    profile = db.relationship('UserProfile', back_populates='user', uselist=False, cascade="all, delete-orphan")
        
    @property
    def password(self):
        return self._password

    @password.setter
    def password(self, raw_password):
        if not raw_password.startswith("pbkdf2:sha256:"):
            self._password = generate_password_hash(raw_password)
        else:
            self._password = raw_password

    def check_password(self, password):
        return check_password_hash(self._password, password)
    
@event.listens_for(User, 'after_insert')
def create_user_profile(mapper, connection, target):
    connection.execute(
        UserProfile.__table__.insert().values(
            user_id=target.id,
            first_name="",
            last_name="",
            bio="",
            profile_picture=""
        )
    )