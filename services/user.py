from werkzeug.security import generate_password_hash, check_password_hash
from models.user import User, db

class UserService:
    def __init__(self):
        self.db = db

    def create_user(self, username, password):
        # Check if the user exists
        existing_user = User.query.filter_by(username=username).first()
        if existing_user:
            raise ValueError("User already exists")
        
        # Hash the password before storing
        hash_password = generate_password_hash(password)

        # Create a new user
        new_user = User(username=username, password=hash_password)
        self.db.session.add(new_user)
        self.db.session.commit()

        return new_user

    def verify_user(self, username, password):
        # Fetch the user by username
        user = User.query.filter_by(username=username).first()
        if not user:
            raise ValueError("Invalid username or password")
        
        # Verify the hashed password
        if not check_password_hash(user.password, password):
            raise ValueError("Invalid username or password")
        
        return user