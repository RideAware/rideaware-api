from flask import Blueprint, request, jsonify, session
from services.UserService.user import UserService

auth_bp = Blueprint("auth", __name__, url_prefix="/api")
user_service = UserService()

@auth_bp.route("/signup", methods=["POST"])
def signup():
    data = request.get_json()
    if not data:
        return jsonify({"message": "No data provided"}), 400
    
    required_fields = ['username', 'password']
    for field in required_fields:
        if not data.get(field):
            return jsonify({"message": f"{field} is required"}), 400
    
    try:
        new_user = user_service.create_user(
            username=data["username"],
            password=data["password"],
            email=data.get("email"),
            first_name=data.get("first_name"),
            last_name=data.get("last_name")
        )
        
        return jsonify({
            "message": "User created successfully", 
            "username": new_user.username,
            "user_id": new_user.id
        }), 201
        
    except ValueError as e:
        return jsonify({"message": str(e)}), 400
    except Exception as e:
        # Log the error
        print(f"Signup error: {e}")
        return jsonify({"message": "Internal server error"}), 500

@auth_bp.route("/login", methods=["POST"])
def login():
    data = request.get_json()
    username = data.get("username")
    password = data.get("password")
    print(f"Login attempt: username={username}, password={password}")
    try:
        user = user_service.verify_user(username, password)
        session["user_id"] = user.id
        return jsonify({"message": "Login successful", "user_id": user.id}), 200
    except ValueError as e:
        print(f"Login failed: {str(e)}")
        return jsonify({"error": str(e)}), 401
    except Exception as e:
        print(f"Login error: {e}")
        return jsonify({"error": "Internal server error"}), 500

@auth_bp.route("/logout", methods=["POST"])
def logout():
    session.clear()
    return jsonify({"message": "Logout successful"}), 200