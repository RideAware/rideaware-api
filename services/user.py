from models.user import User, db
from werkzeug.security import generate_password_hash, check_password_hash

class UserService:
    def create_user(self, username, password):
        if not username or not password:
            return jsonify({"error": "Username and password are required"}), 400
        
        if len(username) < 3 or len(password) < 8:
            return jsonify({"error": "Username must be at least 3 characters and password must be at least 8 characters."}), 400

        
        existing_user = User.query.filter_by(username=username).first()
        if existing_user:
            raise ValueError("User already exists")
        
        hashed_password = generate_password_hash(password)
        new_user = User(username=username, password=hashed_password)
        db.session.add(new_user)
        db.session.commit()
        return new_user

    def verify_user(self, username, password):
        user = User.query.filter_by(username=username).first()
        if not user or not user.check_password(password):
            raise ValueError("Invalid username or password")
        return user
