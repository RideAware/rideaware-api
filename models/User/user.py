from werkzeug.security import generate_password_hash, check_password_hash
from models import db

class User(db.Model):
    __tablename__ = 'users'
    
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    _password = db.Column("password", db.String(255), nullable=False)    
    
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