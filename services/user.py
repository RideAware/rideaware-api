from werkzeug.security import generate_password_hash, check_password_hash
from models.user import User, db

class UserService:
    def create_user(self, username, password):
        existing_user = User.query.filter_by(username=username).first()
        if existing_user:
            raise ValueError("User already exists")
        
        new_user = User(username=username, password=password)
        db.session.add(new_user)
        db.session.commit()
        return new_user

    def verify_user(self, username, password):
        user = User.query.filter_by(username=username).first()
        if not user or not user.check_password(password):
            raise ValueError("Invalid username or password")
        return user
