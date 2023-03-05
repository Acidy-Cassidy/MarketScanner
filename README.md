# MarketScanner


To use:
* Run - Will take a LONG time. It is best to run durning sleep.
  - This takes 45k+ item ids from Jita, pulls the ones that are in the top 1000, highest volume items. 
  - Then it takes this list (which is saved locally next to code) and uses linear regression to show best profit estimates

![ss1](https://user-images.githubusercontent.com/30472756/222958625-972f57dd-7258-4324-b1e5-a32d849d600f.PNG)
  - Once you execute you will see it reteive the item ids from Jita
  - Then it will continue to run for a long while (6 hours) pulling data
 ![ss2](https://user-images.githubusercontent.com/30472756/222958670-b55123e3-c2e8-402c-a9ee-1bfaf539acb0.PNG)
  - These errors are normal, concern is warranted only if 
      *You do not see any progress on the popup console (no traffic) 
      * Or no file is written to your /Documents directory upon completion 


This is a test build used for dev purposes.
