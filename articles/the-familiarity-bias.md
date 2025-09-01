# The familiarity bias

When teams build software, a concept frequently discussed is _simplicity_. The idea that the way we approach problems can be uncomplicated, straightforward.
- _"This code would be much simpler if we used library Y instead"_
- _"We should make this function generic, it will be simpler to reuse"_
- _"Let's use pattern X, it will be simpler"_
- _"Let's remove pattern X, it's too complicated"_

These conversations are reoccuring. Why would anyone think these changes make code simpler? Because the resulting code is less verbose? Because the library is popular? Because we can put a name to the pattern and point to the book where it was first introduced?

The actual reason is obvious but often overlooked: because they have seen a similar code before and people's sense of _what is simple_ is overwhelmingly dominated by familiarity.

This bias is something we as human beings cannot escape. Perception is subjective. We do not have the ability to forget what we know and judge situations objectively.

In this article I want to write about how to identify this bias and defuse it.
But first, let's talk about what simplicity is.

## What is simplicity in software?

Simplicity is hard to define succintly because it spans several dimensions[^1]:
- Something simple does not drown us under its _cognitive load_
- Something simple is _intuitive_ and _self contained_: very little specialised knowledge is required and you can understand it simply by looking at it a few times
- Something simple is _predictable_: once we understand it we are able to predict how it will behave in any given scenario
- Simple things _scale_ and are _composable_: several simple things can be assembled together and remain simple. In contrast, when we assembled complex thing the complexity compounds and explodes, overloading our cognitive abilities for even the smallest of n

## DRAFT NOTES (WIP)

The most common source of complexity in programming is abstraction.




complexity explodes, compounds
m + n or m * n ?
cognitive bounds
knowledge bounds

??
Pushing back on solutions someone perceives as simple can be a complicated conversation. Because in their eyes, you are saying no to simplicity. Questioning  
??

increasing your team's tolerance for complexity?
turning complex into simple?


In this article I want to to try and define simplicity and discuss how to spot & defuse this familiarity bias.
people link abstraction to intelligence. saying no to their abstraction can cause them to feel inadequate


Can feel like a personal attack?

Simplicity = less verbose?
Simplicity = pattern has a name?


debating whether something is simple or not is a fool's errand
it becomes 'personal' -> the guy who doesn't want simple is a bad dev
let's instead shift the conversation to some more concrete attributes of the proposed solutions
blanket word


Does that mean we should get rid of 'patterns' and 'language' and everything that cannot be instantly understood by a toddler?
Admit that there exists situations where simplicity is the what we should optimise for.

the word 'simple' is hiding stuff and reducing quality of conversation / distracting from specific properties of the solution

In the name of simplicity, people push in all kinds of direction (monads? :D)
And saying no is not an option: What kind of incompetent software developer would say not to simplicity?!

We are lacking common definitions.

Everyone thinks what they are doing is simple.
Because it is simple _to them_.

The inability to let go of our own POV is what is blocking us down.

We need forced functions to help us get rid of our own bias.

If its simple to you but not to your teammates -> what knowledge are they missing to turn it into something simple?

What feels 'good' 'bad' -> software instinct? craftsman instinct?

Broadcast what simple is and should be <- Senior IC work?

cognitive ability? knowledge? -> upper bound of simplicity?
expand upper bound by sharing knowledge at the right time

universal simplicity (== regardless of individual?) Does it exist?

Sometimes simplicity is about constraining the complexity and putting it in the right place
The right abstraction will reduce complexity.
The wrong abstraction will increase complexity.
Finding the right abstraction is hard, when we don't know what it is might as well not abstract at all.
You can find an abstraction that feels right _to people who have spent enough time in the problem space_ -> knowledge + familiarity = upper bound


defusing -> 'which of the things you are familiar with are applicable in our case'

you want people to do simple but you don't want them to freeze complexity arises

should we optimise for read or for write ?

[^1]: While writing this article I stumbled upon the excelent talk ["Simple Made Easy"](https://www.youtube.com/watch?v=SxdOUGdseq4) by Rich Hickey. I don't agree with all of the takeaways, but it is definitely the best explanation of what simplicy is I have ever heard. So much so I feel stupid writing this section of the article... Highly recommend watching.
[^2]: If you squint really hard ...
[^3]: Point being: no matter what your personal threshold is it will get exceeded very quickly, because combinatorial logic is explosive
[^4]: The truth is the logic to decide the archetype of a user probably still looks somewhat like implementation 1. We haven't exactly killed the complexity, we just moved it. I swear finding good examples is hard...
