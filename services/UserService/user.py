from models.User.user import User, db

class UserService:
    def create_user(self, username, password):
        if not username or not password:
            raise ValueError("Username and password are required")
        
        if len(username) < 3 or len(password) < 8:
            raise ValueError("Username must be at least 3 characters and password must be at least 8 characters.")
        
        existing_user = User.query.filter_by(username=username).first()
        if existing_user:
            raise ValueError("User already exists")
        
        new_user = User(username=username, password=password)
        db.session.add(new_user)
        db.session.commit()
        return new_user

    def verify_user(self, username, password):
        user = User.query.filter_by(username=username).first()
        if not user:
            print(f"User not found: {username}")
            raise ValueError("Invalid username or password")
        
        if not user.check_password(password):
            raise ValueError("Invalid username or password")
        
        print(f"User verified: {username}")
        return user
