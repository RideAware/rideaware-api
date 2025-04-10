from flask import Blueprint, request, jsonify, session
from services.UserService.user import UserService

auth_bp = Blueprint("auth", __name__, url_prefix="/auth")
user_service = UserService()


@auth_bp.route("/signup", methods=["POST"])
def signup():
    data = request.get_json()
    try:
        new_user = user_service.create_user(data["username"], data["password"])
        return (
            jsonify({"message": "User created successfully", "username": new_user.username}),
            201,
        )
    except ValueError as e:
        return jsonify({"message": str(e)}), 400
    except Exception as e:
        # Log the error
        print(f"Signup error: {e}")
        return jsonify({"message": "Internal server error"}), 500


@auth_bp.route("/login", methods=["POST"])
def login():
    data = request.get_json()
<<<<<<< HEAD
    username = data.get("username")
    password = data.get("password")

    print(f"Login attempt: username={username}, password={password}")

=======
     
>>>>>>> 3ab162d8b88a23ad1d0ef5f72a3162bdd7f75ca8
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
