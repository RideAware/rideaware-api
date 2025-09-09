from models.User.user import User
from models.UserProfile.user_profile import UserProfile
from models import db
import re

class UserService:
    def create_user(self, username, password, email=None, first_name=None, last_name=None):
        if not username or not password:
            raise ValueError("Username and password are required")
        
        if email:
            email_regex = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
            if not re.match(email_regex, email):
                raise ValueError("Invalid email format")
        
        existing_user = User.query.filter(
            (User.username == username) | (User.email == email)
        ).first()
        
        if existing_user:
            if existing_user.username == username:
                raise ValueError("Username already exists")
            else:
                raise ValueError("Email already exists")
        
        if len(password) < 8:
            raise ValueError("Password must be at least 8 characters long")
        
        try:
            new_user = User(
                username=username,
                email=email or "",
                password=password
            )
            
            db.session.add(new_user)
            db.session.flush()
            
            user_profile = UserProfile(
                user_id=new_user.id,
                first_name=first_name or "",
                last_name=last_name or "",
                bio="",
                profile_picture=""
            )
            
            db.session.add(user_profile)
            db.session.commit()
            
            return new_user
            
        except Exception as e:
            db.session.rollback()
            raise Exception(f"Error creating user: {str(e)}")
    
    def verify_user(self, username, password):
        user = User.query.filter_by(username=username).first()
        if not user or not user.check_password(password):
            raise ValueError("Invalid username or password")
        return user