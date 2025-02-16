from models import db
from werkzeug.security import generate_password_hash, check_password_hash

class User(db.Model):
    __tablename__ = 'users'
    
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(80), unique=True, nullable=False)
    password = db.Column(db.String(128), nullable=False)
    
    def __init__(self, username, password, hash_password=True):
        self.username = username
        # Optionally hash the password automatically.
        if hash_password:
            self.password = generate_password_hash(password, method="pbkdf2:sha256")
        else:
            self.password = password

    def check_password(self, password):
        return check_password_hash(self.password, password)
    
    def __repr__(self):
        return f"<User {self.username}>"