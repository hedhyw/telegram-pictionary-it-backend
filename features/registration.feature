Feature: registration
    As a user, I want to be able to join the room
    so that I can participate in the gameplay.

    Scenario: A user joins the room for the first time
        When the user joins the room
        Then they are showed in the list of players
        And they have the correct username
        And they have a zero score
