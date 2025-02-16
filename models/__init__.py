import os
from flask_sqlalchemy import SQLAlchemy
from dotenv import load_dotenv
from urllib.parse import quote_plus

load_dotenv()

PG_USER = quote_plus(os.getenv('PG_USER'))
PG_PASSWORD = quote_plus(os.getenv('PG_PASSWORD'))
PG_HOST = os.getenv('PG_HOST')
PG_PORT = os.getenv('PG_PORT')
PG_DATABASE = os.getenv('PG_DATABASE')

DATABASE_URI = f"postgresql+psycopg2://{PG_USER}:{PG_PASSWORD}@{PG_HOST}:{PG_PORT}/{PG_DATABASE}"

db = SQLAlchemy()

def init_db(app):
    """Initialize the SQLAlchemy app with the configuration."""
    app.config['SQLALCHEMY_DATABASE_URI'] = DATABASE_URI
    app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
    db.init_app(app)
