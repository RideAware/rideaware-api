# routes/auth.py
from flask import Blueprint, request, jsonify, session
from services.UserService.user import UserService

auth_bp = Blueprint('auth', __name__)
user_service = UserService()

@auth_bp.route('/signup', methods=['POST'])
def signup():
    data = request.get_json()
    try:
        new_user = user_service.create_user(data['username'], data['password'])
        return jsonify({"message": "User created successfully", "username": new_user.username}), 201
    except ValueError as e:
        return jsonify({"message": str(e)}), 400

@auth_bp.route('/login', methods=['POST'])
def login():
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')
    
    print(f"Login attempt: username={username}, password={password}")
     
    try:
        user = user_service.verify_user(data['username'], data['password'])
        session['user_id'] = user.id
        return jsonify({"message": "Login successful", "user_id": user.id}), 200
    except ValueError as e:
        print(f"Login failed: {str(e)}")
        return jsonify({"error": str(e)}), 401
    
@auth_bp.route('/logout', methods=['POST'])
def logout():
    session.clear()
    return jsonify({"message": "Logout successful"}), 200