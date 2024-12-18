from flask import Blueprint, request, jsonify
from services.user import UserService

auth_bp = Blueprint('auth', __name__)
user_service = UserService()

@auth_bp.route('/signup', methods=['POST'])
def signup():
    data = request.json
    username = data.get('username')
    password = data.get('password')

    try:
        new_user = user_service.create_user(username, password)
        return jsonify({"message": "User created successfully", "username": new_user.username}), 201
    except ValueError as e:
        return jsonify({"message": str(e)}), 400
    


@auth_bp.route('/login', methods=['POST'])
def login():
    data = request.json
    username = data.get('username')
    password = data.get('password')

    try:
        user = user_service.verify_user(username, password)
        return jsonify({"message": "Login successful", "user_id": user.id}), 200
    except ValueError as e:
        return jsonify({"error": str(e)}), 401