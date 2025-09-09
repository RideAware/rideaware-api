import os
from flask import Flask
from flask_cors import CORS
from dotenv import load_dotenv
from flask_migrate import Migrate
from flask.cli import FlaskGroup

from models import db, init_db
from routes.user_auth import auth

load_dotenv()

app = Flask(__name__)
app.config["SECRET_KEY"] = os.getenv("SECRET_KEY")
app.config["SQLALCHEMY_DATABASE_URI"] = os.getenv("DATABASE")
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False

CORS(app)

init_db(app)
migrate = Migrate(app, db)
app.register_blueprint(auth.auth_bp)


@app.route("/health")
def health_check():
    """Health check endpoint."""
    return "OK", 200

cli = FlaskGroup(app)

if __name__ == "__main__":
    cli()