# Brickwise

Project brief for agents and contributors. This is what the project is for and where it is headed.

## Primary goal

Rebuild an apartments.com-like application to replace managing the caller's properties and tenants, and eventually build leases.

**Features to build (in order of intent):**

1. **Properties** — manage the caller's properties
2. **Tenants** — manage tenants for those properties
3. **Leases** — eventually build leases

## Secondary goal

A Java learning tutorial. The caller is a C#/.NET developer pivoting to Java and needs a pet project for that, so document everything in Java.

## Current branch work

Instructions for in-flight work on this branch (carry as stated):

1. **Fill out the rest of the Java backend.**
2. **Continuously document things for the caller's learning as a C#/.NET developer pivoting to Java,** so that when they review the code they can see how everything works.
3. **Eventually add a new feature,** but the caller will pick what feature to do for that themselves — **do not choose or start a feature now.**

## For a fresh agent

- Optimize for the primary product goal (properties → tenants → leases) unless the task is explicitly about learning/docs.
- Prefer Java for new backend work and for documentation that teaches the stack.
- On this branch: fill out the Java backend and keep learning docs current; wait for the caller before starting any new product feature.
- Do not invent product scope beyond properties, tenants, and leases unless the caller asks for it.

## Layout (this branch)

- `backend/` — Java Spring Boot API
- `frontend/` — Vue 3 app
- `backend/LEARNING.md` — C# /.NET → Java guide for reviewers
- `TENNESSEE_LEASE_REQUIREMENTS.md` — Davidson County / Nashville lease-law notes for a future lease builder (research only; feature not started)
