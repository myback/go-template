from setuptools import setup, find_packages

import go_template


def readme():
    with open('README.md') as f:
        return f.read()


setup(
    name='go-template',
    version=go_template.__version__,
    packages=find_packages(exclude=("tests",)),
    description='Python bindings for go template',
    author='myback',
    long_description=readme(),
    long_description_content_type="text/markdown",
    author_email='54638513+myback@users.noreply.github.com',
    license='MIT',
    url='https://github.com/myback/go-template',
    keywords='golang template go-template bindings wrapper',
    package_data={'go_template': ['bind/template.h', 'bind/template.so']},
    classifiers=['Programming Language :: Python :: 3.5',
                 'Programming Language :: Python :: 3.6',
                 'Programming Language :: Python :: 3.7',
                 'Programming Language :: Python :: 3.8',
                 'Programming Language :: Python :: 3.9',
                 'Programming Language :: Python :: 3.10',
                 'License :: OSI Approved :: MIT License'],
    setup_requires=['wheel']
)
