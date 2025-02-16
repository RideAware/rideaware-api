from flask import Blueprint, request, jsonify
from services.user import UserService

auth_bp = Blueprint('auth', __name__)
user_service = UserService()

@auth_bp.route('/signup', methods=['POST'])
def signup():
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')

    if not username or not password:
        return jsonify({"error": "Username and password are required"}), 400
    
    if len(username) < 3 or len(password) < 8:
        return jsonify({"error": "Username must be at least 3 characters and password must be at least 8 characters."}), 400

    try:
        new_user = user_service.create_user(username, password)
        return jsonify({"message": "User created successfully", "username": new_user.username}), 201
    except ValueError as e:
        return jsonify({"message": str(e)}), 400

@auth_bp.route('/login', methods=['POST'])
def login():
    data = request.get_json()
    username = data.get('username')
    password = data.get('password')

    if not username or not password:
        return jsonify({"error": "Username and password are required"}), 400

    try:
        user = user_service.verify_user(username, password)
        return jsonify({"message": "Login successful", "user_id": user.id}), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 401
