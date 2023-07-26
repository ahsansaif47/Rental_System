import database_utils
import requests
from app import app
from flask import request, jsonify, Response


def check_user_registration(email):
    query = f'''Select email from users where email = {email}'''
    matching_users = database_utils.execute_query(query)
    
    if matching_users:
        return True
            
    return False


@app.route('/register_user', methods=["POST"])
def register_user():
    json_data = request.get_json()
    email = json_data['email']
    password = json_data['password']

    if check_user_registration(email):
        return jsonify({'Status': 'Already registered'})
    else:
        user_insert_query = f'''Insert into users(email, password) Values({email}, {password});'''
        database_utils.execute_query(user_insert_query)
        return Response(status=200)


    
