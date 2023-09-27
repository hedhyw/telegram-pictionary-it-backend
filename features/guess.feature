Feature: guess
    As a player, I want to be able to guess the word
    so that I can win the game.

    Scenario: All players guess the word correctly
        Given there is a started with two players
        When the guesser player guesses the word correctly
        Then the game is finished

    Scenario: Guess the word incorrectly
        Given there is a started game with two players
        When the guesser player guesses the word incorrectly
        Then the game is still in progress

    Scenario: Some players guess the word correct;y
        Given there is a started game with thress players
        When one of the guesser players guesses the word correctly
        Then the game is still in progress
