# 密钥分享 shamir

若有一个富翁，他有千万的身价，他已年老，眼见时日不多，他打算把他毕生的财富
都放到一个数字密码箱里面，等他死后他的三个孩儿们就能打开宝箱，享受幸福人生。
但是这是一个密码箱，是一串3位数的密码，他不能直接告诉他们三个密码是什么，
免得大义灭亲；而且每个孩子都只能知道密码的一部分，只有三个人在一起才能共同打开宝箱。
他该怎么做呢？

## 简单的方法
假如密码箱的密码是123，富翁完全可以让每个孩子记住其中的一位，只有三个孩子都在场，
分别输入自己获得的密码碎片，拼接成一个完整的密码才能打开宝箱，获取富翁的财富。

这个方法虽然能用，但是有点缺陷：
1. 容易被暴力求解，当有两人在场时，完全可以从 0 ~ 9 不断尝试解开密码箱。
2. 密码的位数必须和人数一样（需要保证权利的平等，免得兄弟反目成仇）
3. 密码以明文的方式出现，不安全。

鉴于以上的缺陷，富翁苦思冥想，总于想到了下面的方法...
## 密钥分享（shamir）
富翁随机构造了下面的一个二次方程：
    
    y=x^2+2x+123 (常数 123 就是密码)

然后随机选定下面三个点，交给三个孩子保管：
    
    (1,126),(2,131),(3,138)

任何一个孩子仅凭自己所知道的一个点都无法解开这个未知的一元二次方程，
    
    y=ax^2+bx+c

因为这个方程里面有三个未知数，只能联立三个方程才能解出方程，知道常数c，也就知道密码是什么。
这个方法灵活性比较大，一元二次方程上的点就是一个密钥碎片，可以给多人。方程也不固定是二阶，可兼容多人的情况。
以上只是富翁粗浅理解，如有错误，欢迎指导。


多年以后，富翁死了，他三个孩子相约到宝箱旁，算出了密码，欣喜若狂地打开了宝箱，只见宝箱中写了 24 字：

    富强、民主、文明、和谐、自由、平等、公正、法治、爱国、敬业、诚信、友善。其中“富强、民主、文明、和谐

兄弟们皆泪目，感激父亲地真真教诲。

参考自 : https://www.8btc.com/article/608236 ，