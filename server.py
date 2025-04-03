import os
from flask import Flask
from flask_cors import CORS
from dotenv import load_dotenv

from models import db, init_db
from routes.user_auth import auth

load_dotenv()

app = Flask(__name__)
app.config["SECRET_KEY"] = os.getenv("SECRET_KEY")
app.config["SQLALCHEMY_DATABASE_URI"] = os.getenv("DATABASE")
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

CORS(app)  # Consider specific origins in production

init_db(app)
app.register_blueprint(auth.auth_bp)


@app.route("/health")
def health_check():
    """Health check endpoint."""
    return "OK", 200


with app.app_context():
    db.create_all()

if __name__ == "__main__":
    app.run(debug=True)
