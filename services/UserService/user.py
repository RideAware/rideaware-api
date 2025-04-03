from models.User.user import User, db
import logging

logger = logging.getLogger(__name__)


class UserService:
    def create_user(self, username, password):
        if not username or not password:
            raise ValueError("Username and password are required")

        if len(username) < 3 or len(password) < 8:
            raise ValueError(
                "Username must be at least 3 characters and password must be at least 8 characters."
            )

        existing_user = User.query.filter_by(username=username).first()
        if existing_user:
            raise ValueError("User already exists")

        new_user = User(username=username, password=password)
        db.session.add(new_user)
        try:
            db.session.commit()
        except Exception as e:
            db.session.rollback()
            logger.error(f"Error creating user: {e}")
            raise ValueError("Could not create user") from e
        return new_user

    def verify_user(self, username, password):
        user = User.query.filter_by(username=username).first()
        if not user:
            logger.warning(f"User not found: {username}")
            raise ValueError("Invalid username or password")

        if not user.check_password(password):
            logger.warning(f"Invalid password for user: {username}")
            raise ValueError("Invalid username or password")

        logger.info(f"User verified: {username}")
        return user
