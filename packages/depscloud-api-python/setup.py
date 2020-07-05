from setuptools import setup, find_packages

setup(
    name='depscloud_api',
    version='0.1.7',
    license='MIT',
    description='Python client to the deps.cloud API',
    author='deps.cloud Authors',
    url='https://deps.cloud/',
    install_requires=[
        "protobuf==3.12.2",
        "grpcio==1.30.0"
    ],
    packages=find_packages()
)
