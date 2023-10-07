Feature: start
    As a player, I want to be able to start the game
    so all members can start the gameplay.

    Scenario: Start the game
        Given there are two players in the game
        When the user starts the game
        Then the game is started
        And the leader player receives a random word
