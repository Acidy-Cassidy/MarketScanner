# MarketScanner


To use:
* Run - Will take a LONG time. It is best to run durning sleep.
  - This takes 45k+ item ids from Jita (Data from last 30 days), pulls the ones that are in the top 1000, highest volume items. 
  - Then it takes this list (which is saved locally in your documents directory) which we then use for calculating the profit potential for each item based on the difference between the maximum buy price and minimum sell price, as well as the item's volume. The profit margin is calculated as (sell price - buy price) / sell price, and then multiplied by the item's volume to obtain a rough estimate of the potential profit. This calculation does not take into account transaction fees or other expenses, so it is important to use this estimate as a starting point and to conduct further research and analysis before making any investment decisions.

![ss1](https://user-images.githubusercontent.com/30472756/222958625-972f57dd-7258-4324-b1e5-a32d849d600f.PNG)
  - Once you execute you will see it retrieve the item ids from Jita
  - Then it will continue to run for a long while (6 hours) pulling data
 ![ss2](https://user-images.githubusercontent.com/30472756/222958670-b55123e3-c2e8-402c-a9ee-1bfaf539acb0.PNG)
  - These errors are normal, concern is warranted only if 
      *You do not see any progress on the popup console (no traffic) 
      * Or no file is written to your /Documents directory upon completion 


This is a test build used for dev purposes.
