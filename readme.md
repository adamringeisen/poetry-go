# Poetry Generator 2.0

## 10-12-23
Rewritten from python. Only minor differences. 

1. The lines for the poem are now embedded in the executable. It's like 100k and considering the size of golang execs, I don't thing it'll be an issue. Functionally this means that it might be possible to two (or more) lines to be the same because now I am simply generating four random numbers and I think it's possible to the generator to return the same number in succession whereas before I was using sqlite to get four distinct values from the table at random and this wasn't possible. I think it's pratically impossible and poems have lines that repeat all the time so I think it will be fine.

2. The colors around each poem are now generated a little differently because I was trying to weed out the eyebleeding colors that you get with just three random bytes as RGB. The colors are now just a struct of three bytes and I add a (different) random number between 1 and 100 to each one. This keeps them sort of sane. Lots of earth tones so it's an improvement over the previous pure random colors.

3. The performance is insanely better. It's on another order of magnitude. The python version started dropping requests at just 300/s while this version his 2000/s with no issues.
