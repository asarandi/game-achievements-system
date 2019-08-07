package main

type AchievementSlugFunction struct {achievement Achievement; slug string; function func(Stat)bool}

//the string should match the `slug` in achievements sql table
var ASF = []AchievementSlugFunction {
    {achievement: Achievement{}, slug: "sharpshooter", function: isSharpshooterAward},
    {achievement: Achievement{}, slug: "bruiser",      function: isBruiserAward},
    {achievement: Achievement{}, slug: "veteran",      function: isVeteranAward},
    {achievement: Achievement{}, slug: "bigwinner",    function: isBigwinnerAward},
}

// “sharpshooter award” – a user receives this for landing 75% of their attacks, assuming they have at least attacked once.
func isSharpshooterAward(stat Stat) bool {
    if (stat.NumAttacks > 0) {
        if (stat.NumHits * 100 / stat.NumAttacks) >= 75 {
            return true
        }
    }
    return false
}

//“bruiser” award – a user receives this for doing more than 500 points of damage in one game
func isBruiserAward(stat Stat) bool {
    if (stat.AmountDamage >= 500) {
        return true
    }
    return false
}

//“veteran” award – a user receives this for playing more than 1000 games in their lifetime.
func isVeteranAward(stat Stat) bool {
    var count int
    if  DB.Find(&Stat{}, Stat{MemberID: stat.MemberID}).Count(&count).Error == nil {
       if count >= 1000 {
            return true
       }
    }
    return false
}

//“big winner” award – a user receives this for having 200 wins
func isBigwinnerAward(stat Stat) bool {
    var count int
    if  DB.Find(&Stat{}, Stat{MemberID: stat.MemberID, IsWinner: true}).Count(&count).Error == nil {
       if count >= 200 {
            return true
       }
    }
    return false
}
