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
    return (stat.NumAttacks > 0) && ((stat.NumHits * 100 / stat.NumAttacks) >= 75)
}

//“bruiser” award – a user receives this for doing more than 500 points of damage in one game
func isBruiserAward(stat Stat) bool {
    return stat.AmountDamage >= 500
}

//“veteran” award – a user receives this for playing more than 1000 games in their lifetime.
func isVeteranAward(stat Stat) bool {
    var count int
    DB.Find(&Stat{}, Stat{MemberID: stat.MemberID}).Count(&count)
    return count >= 1000
}

//“big winner” award – a user receives this for having 200 wins
func isBigwinnerAward(stat Stat) bool {
    var count int
    DB.Find(&Stat{}, Stat{MemberID: stat.MemberID, IsWinner: true}).Count(&count)
    return count >= 200
}
