VENV := .venv
PYTHON := $(VENV)/bin/python

.PHONY: setup test lint format pre-commit docs coverage serve-coverage clean help all

install-uv:
	@echo "Installing UV package manager..."
	@curl -LsSf https://astral.sh/uv/install.sh | sh
	@echo "UV installed successfully!"

setup:
	@if ! command -v uv >/dev/null 2>&1; then \
		${MAKE} install-uv; \
	fi
	@echo "Setting up development environment..."
	@uv venv $(VENV)
	@uv pip install --python $(PYTHON) -e ".[dev]"
	@echo "Development environment setup complete!"

install:
	@uv pip install --python $(PYTHON) -e .

lint:
	@echo "Running linters..."
	@uv run --python $(PYTHON) ruff check .
	@uv run --python $(PYTHON) mypy --show-error-codes enterprise_ai/

format:
	@echo "Formatting code..."
	@uv run --python $(PYTHON) ruff format .
	@uv run --python $(PYTHON) ruff check --fix .
	@echo "Formatting Markdown files..."
	@which mdformat >/dev/null 2>&1 || uv pip install mdformat
	@mdformat .

pre-commit:
	@echo "Running pre-commit hooks..."
	@uv run --python $(PYTHON) pre-commit run --all-files

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf build/ dist/ *.egg-info/ .pytest_cache/ .ruff_cache/ .mypy_cache/ htmlcov/ .coverage logs/ workspace/
	@find . -type d -name __pycache__ -exec rm -rf {} +

notebook:
	@echo "Checking for Jupyter Notebook..."
	@if ! $(VENV)/bin/pip show notebook > /dev/null 2>&1; then \
		echo "Installing Jupyter Notebook..."; \
		uv pip install --python $(PYTHON) notebook; \
	fi
	@echo "Launching Jupyter Notebook..."
	@$(VENV)/bin/jupyter notebook

all: lint pre-commit

help:
	@echo "Enterprise-AI Development Makefile"
	@echo "=================================="
	@echo "setup          - Create virtual env and install deps"
	@echo "install        - Install package in dev mode"
	@echo "lint           - Run static analysis checks"
	@echo "format         - Format and fix code"
	@echo "clean          - Remove build artifacts"
	@echo "pre-commit     - Run all pre-commit checks"
	@echo "notebook       - Launch Jupyter Notebook"
	@echo "all            - Run full quality checks (lint + test + coverage)"