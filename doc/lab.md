---
title: LAB
section: 1
header: User Manual
footer: lab 1.0.0
date: January 14, 2024
---
# NAME
lab - A program for conducting CIT labs

# SYNOPSIS
**lab** SUB-COMMAND [*args*]

SUB-COMMAND := {start | check | status | current | flag | submit | score | help}

# DESCRIPTION
**lab** is a CIT lab system, built to guide students through lab assignments, providing feedback and tips, tracking progress and giving students additional experience using the Unix command line. 

Each lab consists of multiple **steps**, and each step composed of one or more **tasks**. Tasks will be either a command that must be typed exactly as shown, or a command that is described by that students must figure out. After completing all of the tasks in the step, students can check to see if they have performed them properly, and to try again if not. Once a step has been checked successfully, students may move on to the next step. After all steps are completed, the student will submit the lab for grading. This will provide a **submission code** that the student will submit to Brightpsace. 

Labs may also contain **flags**. A **flag** is a four digit number that appears after the text "flag_". For example: "flag_1234". These will appear when running certain commands. Students may "capture the flag" by using a command. Some flags will appear as long as the student executes the steps properly. Other flags, called **bonus flags**, are hidden and require the student to explore and run commands that are not part of the lab. These are worth bonus points. 

**lab** also provides sub-commands for viewing the steps of the lab, the details of the current step

# SCORING

The score is calculated by adding the number of successfully completed steps to the number of captured flags and the number of bonus flags. The highest possible score is the number of steps plus the number of flags. This means that if you find bonus flags, your score could potentially be higher than the total possible. 

For example, if a lab has 7 steps, has 3 flags and 2 bonus flags, the total possible is:

    7 + 3 => 10 total possible points

If you capture the 2 bonus flags, your score can be 12 out of 10. Those two points will be bonus points toward your final grade in the class. 


# SUB-COMMANDS
**help** 
: display help message

**start** <*lab-id*>
: starts a new lab. When you run the start command, you must provide the lab id (see EXAMPLES), such as "2-1", which will be given to you in the assignment description in Brightspace. Only one lab may be in progress at a time. If you try to start a lab before submitting one that is in progress, you will get a notice and the opportunity to submit the previous lab. 

**steps** 
: view a list of steps for the currently active lab.

**current**
: view the details and tasks of the current step.

**check**
: check to see if you have completed the current step properly. If not, you may try again and continue checking until you succeed. Upon success, this command will prompt you to move to the next step if you are ready. 

**flag** <*flag-number*>
: captures a flag or bonus flag. This command requires a four digit flag number argument (see EXAMPLES) that will be seen in the output of a command during the lab (regular flags) or while exploring on your own (bonus flags). This can only capture flags that are valid for the active lab. 

**submit**
: submit the lab for grading. You can submit a lab without completing all of the steps or capturing all regular flags. Your score report will be printed to the screen, along with a unique **submission code** that you will submit to Brightspace to verify that you have completed the lab.

**score**
: review the score and results of a submitted lab. This will display the number of completed steps, captured flags and bonus flags and total score. 

# EXAMPLES
**lab start 5-2**
: Begins lab number 5-2. 

**lab flag 1234**
: Captures flag number 1234, if it is a valid flag for the current lab. 

# AUTHORS
Written by Jeremy Doolin <jdoolin@wvncc.edu>.

# BUGS
Submit bug reports to Jeremy Doolin <jdoolin@wvncc.edu>