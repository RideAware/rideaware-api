import os
from flask import Flask
from models import db, init_db 
from routes.user_auth.auth import auth_bp
from dotenv import load_dotenv
from flask_cors import CORS

load_dotenv()

app = Flask(__name__)
CORS(app)

app.config['SQLALCHEMY_DATABASE_URI'] = os.getenv('DATABASE')
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

init_db(app)

app.register_blueprint(auth_bp)

with app.app_context():
    db.create_all()

if __name__ == '__main__':
    app.run(debug=True)
