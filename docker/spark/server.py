import os
from flask import Flask
from flask_restful import reqparse, abort, Api, Resource
import subprocess

if not os.path.exists('/tmp/spark-events'):
    os.makedirs('/tmp/spark-events')

app = Flask(__name__)
api = Api(app)

parser = reqparse.RequestParser()

parser.add_argument('input')
parser.add_argument('output')
parser.add_argument('old')

spark_path = f'{os.getcwd()}/sparkjob.jar'
print(spark_path)

class SubmitJob(Resource):
    def post(self):
        args = parser.parse_args()
        print(args)
        
        if args['old'] is None:
            subprocess.run(['/opt/spark/bin/spark-submit', spark_path, args['input'], args['output']])
        else:
            subprocess.run(['/opt/spark/bin/spark-submit', spark_path, args['input'], args['output'], args['old']])
        return {'status': 'ok'}, 200


api.add_resource(SubmitJob, '/submit')

if __name__ == '__main__':
    app.run(debug=True)